import React, { useState } from "react";

export function Body() {
  const [count, setCount] = useState(0);
  return (
    <div className="flex w-full flex-col items-center justify-center">
      <h1>This is client-side content from React!</h1>
      <p>Counter: {count}</p>
      <button onClick={() => setCount(count + 1)}>Increment</button>
    </div>
  );
}
