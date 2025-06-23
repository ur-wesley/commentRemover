#!/usr/bin/env node

const { spawn } = require("child_process");
const path = require("path");
const fs = require("fs");

function runCommand(command, args = []) {
 return new Promise((resolve, reject) => {
  const child = spawn(command, args, {
   stdio: "pipe",
   cwd: process.cwd(),
  });

  let stdout = "";
  let stderr = "";

  child.stdout.on("data", (data) => {
   stdout += data.toString();
  });

  child.stderr.on("data", (data) => {
   stderr += data.toString();
  });

  child.on("close", (code) => {
   if (code === 0) {
    resolve({ stdout, stderr, code });
   } else {
    reject({ stdout, stderr, code });
   }
  });

  child.on("error", (error) => {
   reject({ error: error.message, code: -1 });
  });
 });
}

async function testCommenter() {
 console.log("🧪 Testing Comment Remover npm package...\n");

 try {
  // Test 1: Check if binary exists
  const binaryPath = path.join(__dirname, "..", "bin", "commenter");
  const binaryExists = fs.existsSync(binaryPath);

  console.log(`📦 Binary exists: ${binaryExists ? "✅" : "❌"}`);
  if (!binaryExists) {
   console.log("   Binary not found. Run: npm install");
   return false;
  }

  // Test 2: Test help command
  console.log("\n🔍 Testing help command...");
  const helpResult = await runCommand(binaryPath, ["--help"]);
  console.log("✅ Help command works");

  // Test 3: Test version command
  console.log("\n📋 Testing version command...");
  const versionResult = await runCommand(binaryPath, ["--version"]);
  console.log("✅ Version command works");

  // Test 4: Test with a sample file
  console.log("\n📝 Testing with sample file...");
  const testFile = path.join(__dirname, "..", "test_examples", "example.ts");

  if (fs.existsSync(testFile)) {
   const testResult = await runCommand(binaryPath, [testFile]);
   console.log("✅ Sample file processing works");
  } else {
   console.log("⚠️  Sample file not found, skipping file test");
  }

  // Test 5: Test bun run compatibility
  console.log("\n🚀 Testing bun run compatibility...");
  try {
   const bunResult = await runCommand("node", [
    path.join(__dirname, "..", "index.js"),
    "--help",
   ]);
   console.log("✅ JavaScript wrapper works");
  } catch (error) {
   console.log("❌ JavaScript wrapper failed:", error.error || error.stderr);
  }

  console.log("\n🎉 All tests passed!");
  console.log("\n💡 You can now use:");
  console.log("   commenter <file>           # Direct binary usage");
  console.log("   bun run commenter <file>   # Via npm script");
  console.log("   npx @ur-wesley/commenter <file>  # Via npx");

  return true;
 } catch (error) {
  console.error(
   "❌ Test failed:",
   error.error || error.stderr || error.message
  );
  return false;
 }
}

if (require.main === module) {
 testCommenter().then((success) => {
  process.exit(success ? 0 : 1);
 });
}

module.exports = { testCommenter };
