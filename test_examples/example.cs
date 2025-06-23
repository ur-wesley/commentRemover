using System;

// This is a single-line comment
/* This is a multi-line comment
   that spans multiple lines */

namespace Example
{
    // Class comment
    public class Program
    {
        // Property comment
        public string Name { get; set; }
        
        /* Single-line multi-line comment */
        
        // Constructor comment
        public Program()
        {
            Console.WriteLine("Hello World"); // Inline comment
        }
        
        /* 
         * Multi-line comment with asterisks
         * // This should NOT be removed
         * End of comment
         */
        
        public void Test()
        {
            string message = "String with // comment inside";
            Console.WriteLine(message);
        }
    }
} 