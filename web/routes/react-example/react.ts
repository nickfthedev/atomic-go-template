// This is the entry point for the react components
// The build script will search for react.ts and react.js files and bundle them into a single file
// The out.js can then be included in the templ file
// 

import React from 'react';
import { createRoot } from 'react-dom/client';
import { Body } from './react-component';

// Add the body React component.
const contentRoot = document.getElementById('react-content');
if (!contentRoot) {
  throw new Error('Could not find element with id react-content');
}

const contentReactRoot = createRoot(contentRoot);
contentReactRoot.render(React.createElement(Body));