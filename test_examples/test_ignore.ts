// This is a regular comment that should be removed
// @ts-ignore This comment should be preserved
// @deprecated This function is deprecated
// TODO: This should be preserved
// FIXME: This should also be preserved

function example() {
 // Regular inline comment
 const x = 42; // @ts-ignore inline ignore comment

 /* Regular multi-line comment */

 /* @ts-ignore multi-line ignore comment */

 return x;
}

// Another regular comment
// @ts-expect-error This should be preserved too
