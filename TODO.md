

1. Add normal error handling to the ServeBot function. - 10% done

2. Delete non-needed debug lines - done.

3. Add the substitute functions - done.

4. Instead of logging errors into the console, we could try setting up an email and sending them directly there - that way, we will be able to get on top of it faster.

5. Standardize error handling.

6. The guide for users should not be on our github repository, seeing as it also stores out bottoken.txt and the password for accessing the database.

7. Make actually readable commands - by storing the current moment at which we are in the chat.

8. Adding exceptions to both banned words (done for substitute words).

9. Rewrite all database functions to standard (error, string) format to make them easier to handle. - done, also deleted non-needed functions.

10. Add all of the possible commands to the Telegram menu

11. At some point in the future, reorganize the code to make it actually readable.

12. Make the anti-spam check a function instead of a pile of code inside the handler.

13. Implement at least basic protection against SQL injections.


