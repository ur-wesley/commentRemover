SELECT * FROM users WHERE active = 1;
/* Multi-line SQL comment
   should be preserved */
INSERT INTO logs (message) VALUES ('String with -- fake comment');
COMMIT; 
