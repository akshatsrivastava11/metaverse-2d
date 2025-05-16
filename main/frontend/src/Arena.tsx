//@ts-nocheck
import React, { use, useEffect, useRef, useState } from 'react';

const Arena = () => {
  const canvasRef = useRef<any>(null);
  const wsRef = useRef<any>(null);
  const [currentUser, setCurrentUser] = useState<any>({});
  const [users, setUsers] = useState(new Map());
  const [params, setParams] = useState({ token: '', spaceId: '' });
  const currentUserRef=useRef(null)
  useEffect(() => {
  currentUserRef.current = currentUser;
}, [currentUser]);

  // Initialize WebSocket connection and handle URL params
  useEffect(() => {
    const urlParams = new URLSearchParams(window.location.search);
    console.log(urlParams)
    const token = urlParams.get('token') || '';
    const spaceId = urlParams.get('spaceId') || '';
    setParams({ token, spaceId });
    console.log("token and space id is",{token,spaceId})

    // Initialize WebSocket
    wsRef.current = new WebSocket('ws://localhost:3002/'); // Replace with your WS_URL
    
    wsRef.current.onopen = () => {
      // Join the space once connected
      wsRef.current.send(JSON.stringify({
        type: 'join',
        payload: {
          spaceId,
          token
        }
      }));
    };
    // console.log("a ")

    wsRef.current.onmessage = (event: any) => {
      // console.log(even)
      console.log("event.data is ",event.data)
      const message = JSON.parse(event.data);
      console.log("jsonparsed event.data is ",message)
      handleWebSocketMessage(message);
    };

    return () => {
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, []);

  const handleWebSocketMessage = (message: any) => {
    switch (message.type) {
            case 'welcome':
      // Set current user and existing users
      setCurrentUser({
        userId: message.payload.userId,
        x: message.payload.x,
        y: message.payload.y,
      });
      // Populate users map with existing users (excluding yourself)
      setUsers(prev => {
        const newUsers = new Map(prev);
        message.payload.users.forEach((user: any) => {
          if (user.userId !== message.payload.userId) {
            newUsers.set(user.userId, user);
          }
        });
        return newUsers;
      });
      break;
      case 'user-joined':
        console.log("A user joined event occured")  
    setUsers(prev => {
  const newUsers = new Map(prev);
  newUsers.set(message.payload.userId, {
    x: message.payload.x,
    y: message.payload.y,
    userId: message.payload.userId
    });
    console.log(newUsers)
    return newUsers;
    
    });

        break;

case 'movement':
  setUsers(prev => {
    const newUsers = new Map(prev);
    newUsers.set(message.payload.userId, {
      ...newUsers.get(message.payload.userId),
      x: message.payload.x,
      y: message.payload.y,
      userId: message.payload.userId,
    });
    return newUsers;
  });
  // ADD THIS BLOCK:
  if (message.payload.userId === currentUserRef.current.userId) {
    setCurrentUser(prev => ({
      ...prev,
      x: message.payload.x,
      y: message.payload.y,
    }));
  }
  break;
      case 'user-left':
        setUsers(prev => {
          const newUsers = new Map(prev);
          newUsers.delete(message.payload.userId);
          return newUsers;
        });
        break;
    }
  };

  // Handle user movement
  const handleMove = (newX: any, newY: any) => {
    if (!currentUser) return;
    console.log("A handle movement event occured")
    console.log("The newcordinates are",XMLDocument,newX,newY)

    // Send movement request
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN){
    wsRef.current.send(JSON.stringify({
      type: 'move',
      payload: {
        x: newX,
        y: newY,
        userId: currentUserRef.current.userId
      }
    }));
  };
}

  // Draw the arena
useEffect(() => {
  console.log("in the making of the canvas")
  const canvas = canvasRef.current;
  if (!canvas) return;
  const ctx = canvas.getContext('2d');
  ctx.clearRect(0, 0, canvas.width, canvas.height);

  // Draw grid
  ctx.strokeStyle = '#eee';
  for (let i = 0; i < canvas.width; i += 50) {
    ctx.beginPath();
    ctx.moveTo(i, 0);
    ctx.lineTo(i, canvas.height);
    ctx.stroke();
  }
  for (let i = 0; i < canvas.height; i += 50) {
    ctx.beginPath();
    ctx.moveTo(0, i);
    ctx.lineTo(canvas.width, i);
    ctx.stroke();
  }

  // Draw current user (red)
  if (currentUser && currentUserRef.current.x !== undefined && currentUserRef.current.y !== undefined)  {
    ctx.beginPath();
    ctx.arc(currentUserRef.current.x, currentUserRef.current.y, 20, 0, Math.PI * 2);
    ctx.fillStyle = '#FF6B6B';
    ctx.fill();
  }

  // Draw other users (blue)
  users.forEach(user => {
    if (user.userId === currentUserRef.current.userId) return; // skip yourself
    if (user.x === undefined || user.y === undefined) return;
    ctx.beginPath();
    ctx.arc(user.x, user.y, 20, 0, Math.PI * 2);
    ctx.fillStyle = '#4ECDC4';
    ctx.fill();
  });
}, [currentUserRef.current, users]);


  const handleKeyDown = (e: any) => {
    if (!currentUserRef.current) return;
    // console.log("An event occured")
    const { x, y } = currentUserRef.current;
    switch (e.key) {
      case 'ArrowUp':
        console.log("Arrowup mann")
        console.log("previous value of xand y",x,y)
        handleMove(x, y - 1);
        break;
      case 'ArrowDown':
       console.log("ArrowDown mann")

        handleMove(x, y + 1);
        break;
      case 'ArrowLeft':
                console.log("ArrowLeft mann")

        handleMove(x - 1, y);
        break;
      case 'ArrowRight':
                console.log("ArrowRight mann")

        handleMove(x + 1, y);
        break;
    }
  };

  return (
    <div className="p-4" onKeyDown={handleKeyDown} tabIndex={0}>
        <h1 className="text-2xl font-bold mb-4">Arena</h1>
        <div className="mb-4">
          <p className="text-sm text-gray-600">Token: {params.token}</p>
          <p className="text-sm text-gray-600">Space ID: {params.spaceId}</p>
          <p className="text-sm text-gray-600">Connected Users: {users.size + (currentUserRef.current ? 1 : 0)}</p>
        </div>
        <div className="border rounded-lg overflow-hidden">
          <canvas
            ref={canvasRef}
            width={2000}
            height={2000}
            className="bg-white"
          />
        </div>
        <p className="mt-2 text-sm text-gray-500">Use arrow keys to move your avatar</p>
    </div>
  );
};

export default Arena;