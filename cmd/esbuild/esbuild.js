// Searches the directory for files with the name react.ts or react.js or *.react.ts or *.react.js and bundles them into a single file
// excluding out.js files
// creates a out.js file in every directory or *.out.js files if there are multiple .react.ts or .react.js files

import fs from "fs";
import path from "path";
import esbuild from "esbuild";

const routesDir = path.join("web", "routes");
const componentsDir = path.join("web", "components");

// Iterate over all directories and recursively iterate over all subdirectories
function iterateOverDirectories(dir) {
  const files = fs.readdirSync(dir);
  for (const file of files) {
    const filePath = path.join(dir, file);
    if (fs.statSync(filePath).isDirectory()) {
      iterateOverDirectories(filePath);
    }
    // check for react.ts, react.js
    // You could also bundle multiple files they need to end with .react.ts, .react.js
    if (file.endsWith("react.ts") || file.endsWith("react.js")) {
      // Check if this file is a file with *.react.ts or *.react.js
      if (file.endsWith(".react.ts") || file.endsWith(".react.js")) {
        // Remove the .react from the file name and the .ts or .js from the end
        const subFileName = file
          .replace(".react", "")
          .replace(".ts", "")
          .replace(".js", "");
        bundleFile(filePath, subFileName);
      } else {
        bundleFile(filePath);
      }
    }
  }
}

// Bundles the file into a single file
// If subFileName is provided, it will be used as the name of the file
function bundleFile(file, subFileName) {
  let filename = "out.js";
  if (subFileName) {
    filename = subFileName + ".out.js";
  }
  esbuild.buildSync({
    entryPoints: [file],
    outfile: path.join(
      "web",
      "embed",
      "assets",
      "react",
      path.dirname(file),
      filename
    ),
    bundle: true,
    minify: true,
    sourcemap: true,
    loader: { ".js": "jsx", ".ts": "tsx" },
    format: "esm",
    target: ["es6"],
    external: [], // This ensures all dependencies are bundled
    define: {
      "process.env.NODE_ENV": '"production"',
    },
  });
  console.log(`Bundled ${file} into ${filename} in ${path.dirname(file)}`);
}

// Delete all files and directories in a directory
function deleteFilesInDirectory(dir) {
  fs.rmSync(dir, { recursive: true, force: true });
}

// Delete all files and directories in the web/embed/assets/js/esbuild directory
deleteFilesInDirectory(path.join("web", "embed", "assets", "react"));

// Iterate over all directories and recursively iterate over all subdirectories
// This calls the bundleFile function for each file
iterateOverDirectories(routesDir);
iterateOverDirectories(componentsDir);
console.log("Bundling complete");
