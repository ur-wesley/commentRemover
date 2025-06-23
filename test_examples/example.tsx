import React from "react";

const MyComponent: React.FC = () => {
 /* This is a regular multi-line comment
     that spans multiple lines */

 const message = "Hello World";

 return (
  <div>
   {/* This is a JSX comment that should be detected */}
   <h1>{message}</h1>

   {/* 
        This is a multi-line JSX comment
        that spans multiple lines
      */}

   <p>This is some text</p>

   {/* Another JSX comment */}

   {/* 
        JSX comments can contain special characters:
        - Dashes: --
        - Slashes: //
        - Asterisks: ***
      */}
  </div>
 );
};

export default MyComponent;
