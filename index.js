#!/usr/bin/env node

const { spawn } = require("child_process");
const path = require("path");
const fs = require("fs");
const os = require("os");

function getBinaryPath() {
 const platform = os.platform();
 const extension = platform === "win32" ? ".exe" : "";
 return path.join(__dirname, "bin", `commenter${extension}`);
}

function main() {
 const binaryPath = getBinaryPath();

 if (!fs.existsSync(binaryPath)) {
  console.error("❌ Commenter binary not found. Please run: npm install");
  console.error("   Or if you have Go installed: go build -o bin/commenter");
  process.exit(1);
 }

 const args = process.argv.slice(2);

 const child = spawn(binaryPath, args, {
  stdio: "inherit",
  cwd: process.cwd(),
 });

 child.on("error", (error) => {
  console.error("❌ Failed to execute commenter:", error.message);
  process.exit(1);
 });

 child.on("close", (code) => {
  process.exit(code);
 });
}

if (require.main === module) {
 main();
}

module.exports = { main, getBinaryPath };
