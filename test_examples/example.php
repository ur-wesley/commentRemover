<?php
// This is a single-line comment
/* This is a multi-line comment
   that spans multiple lines */

class Example {
    // Class property comment
    private $name = "test";
    
    /* Single-line multi-line comment */
    
    public function __construct() {
        // Constructor comment
        echo "Hello World"; // Inline comment
    }
    
    /* 
     * Multi-line comment with asterisks
     * // This should NOT be removed
     * End of comment
     */
    
    public function test() {
        $message = "String with // comment inside";
        return $message;
    }
}
?> 