#!/usr/bin/env node

const { execSync } = require("child_process");
const fs = require("fs");
const path = require("path");
const os = require("os");

const BINARY_NAME =
 process.platform === "win32" ? ".\\commenter-test.exe" : "./commenter-test";

console.log("🧪 Running npm package tests...\n");

function testBinaryExists() {
 console.log("📦 Test 1: Check binary existence...");
 try {
  if (!fs.existsSync(BINARY_NAME.replace("./", "").replace(".\\", ""))) {
   throw new Error(`Binary ${BINARY_NAME} not found`);
  }
  console.log("✅ Binary exists\n");
  return true;
 } catch (error) {
  console.error("❌ Binary test failed:", error.message);
  return false;
 }
}

function testHelpCommand() {
 console.log("📋 Test 2: Help command...");
 try {
  const output = execSync(`${BINARY_NAME} -h`, {
   encoding: "utf8",
   timeout: 5000,
  });

  const expectedStrings = [
   "Comment Remover",
   "USAGE:",
   "OPTIONS:",
   "--write",
   "--recursive",
   "--no-color",
   "--help",
  ];

  for (const expected of expectedStrings) {
   if (!output.includes(expected)) {
    throw new Error(`Help output missing expected string: "${expected}"`);
   }
  }

  console.log("✅ Help command works correctly\n");
  return true;
 } catch (error) {
  console.error("❌ Help command test failed:", error.message);
  return false;
 }
}

function testFileProcessing() {
 console.log("📄 Test 3: File processing...");

 const testDir = fs.mkdtempSync(path.join(os.tmpdir(), "commenter-test-"));
 const testFile = path.join(testDir, "test.go");

 try {
  const testContent = `package main

import "fmt"

func main() {
    fmt.Println("Hello")
    /* This should stay
       // This should also stay
    */
    fmt.Println("Done")
}
`;

  fs.writeFileSync(testFile, testContent);

  const output = execSync(`${BINARY_NAME} "${testFile}"`, {
   encoding: "utf8",
   timeout: 10000,
  });

  if (!output.includes("Comments removed: 2")) {
   if (!output.includes("Comments removed:")) {
    throw new Error("Expected comments to be removed");
   }
  }

  if (
   !output.includes("Original lines:") ||
   !output.includes("Remaining lines:")
  ) {
   throw new Error("Missing expected statistics in output");
  }

  console.log("✅ File processing works correctly\n");
  return true;
 } catch (error) {
  console.error("❌ File processing test failed:", error.message);
  return false;
 } finally {
  try {
   fs.rmSync(testDir, { recursive: true, force: true });
  } catch (cleanupError) {
   console.warn("⚠️  Cleanup warning:", cleanupError.message);
  }
 }
}

function testDirectoryProcessing() {
 console.log("📁 Test 4: Directory processing...");

 const testDir = fs.mkdtempSync(path.join(os.tmpdir(), "commenter-dir-test-"));

 try {
  const file1 = path.join(testDir, "file1.go");
  const file2 = path.join(testDir, "file2.js");

  fs.writeFileSync(file1, "package main\n// comment\nfunc main() {}");
  fs.writeFileSync(file2, '// comment\nconsole.log("hello");');

  const output = execSync(`${BINARY_NAME} "${testDir}"`, {
   encoding: "utf8",
   timeout: 10000,
  });

  if (!output.includes("Batch Processing Summary:")) {
   throw new Error("Missing batch processing summary");
  }

  if (!output.includes("Files processed:")) {
   throw new Error("Expected files to be processed");
  }

  console.log("✅ Directory processing works correctly\n");
  return true;
 } catch (error) {
  console.error("❌ Directory processing test failed:", error.message);
  return false;
 } finally {
  try {
   fs.rmSync(testDir, { recursive: true, force: true });
  } catch (cleanupError) {
   console.warn("⚠️  Cleanup warning:", cleanupError.message);
  }
 }
}

function testErrorHandling() {
 console.log("❗ Test 5: Error handling...");

 try {
  try {
   execSync(`${BINARY_NAME} /non/existent/file.go`, {
    encoding: "utf8",
    timeout: 5000,
    stdio: "pipe",
   });
   throw new Error("Expected command to fail with non-existent file");
  } catch (cmdError) {
   if (cmdError.status === 0) {
    throw new Error("Command should have failed but succeeded");
   }
  }

  const testDir = fs.mkdtempSync(
   path.join(os.tmpdir(), "commenter-error-test-")
  );
  const unsupportedFile = path.join(testDir, "test.py");

  try {
   fs.writeFileSync(unsupportedFile, '# Python comment\nprint("hello")');

   execSync(`${BINARY_NAME} "${unsupportedFile}"`, {
    encoding: "utf8",
    timeout: 5000,
    stdio: "pipe",
   });
   throw new Error("Expected command to fail with unsupported file type");
  } catch (cmdError) {
   if (cmdError.status === 0) {
    throw new Error("Command should have failed but succeeded");
   }
  } finally {
   fs.rmSync(testDir, { recursive: true, force: true });
  }

  console.log("✅ Error handling works correctly\n");
  return true;
 } catch (error) {
  console.error("❌ Error handling test failed:", error.message);
  return false;
 }
}

async function runAllTests() {
 const tests = [
  testBinaryExists,
  testHelpCommand,
  testFileProcessing,
  testDirectoryProcessing,
  testErrorHandling,
 ];

 let passed = 0;
 let failed = 0;

 for (const test of tests) {
  try {
   if (test()) {
    passed++;
   } else {
    failed++;
   }
  } catch (error) {
   console.error("❌ Test execution error:", error.message);
   failed++;
  }
 }

 console.log("📊 Test Results:");
 console.log(`✅ Passed: ${passed}`);
 console.log(`❌ Failed: ${failed}`);
 console.log(
  `📈 Success Rate: ${Math.round((passed / (passed + failed)) * 100)}%\n`
 );

 if (failed > 0) {
  console.log(
   "🚨 Some tests failed. Please check the output above for details."
  );
  process.exit(1);
 } else {
  console.log("🎉 All tests passed! The npm package is working correctly.");
  process.exit(0);
 }
}

runAllTests().catch((error) => {
 console.error("💥 Test runner error:", error);
 process.exit(1);
});
