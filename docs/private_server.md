# Make Lenpaste server private
You can password protect the server.
In this case, everyone will be able to view the pastes, and only authorized users will be able to create them.
Authorization is done using the HTTP Basic Authentication standard.

Create a lenpasswd file.
If you are using Docker, put it under `/data/lenpasswd`.
If you just run Lenpaste, specify the path to the file using the `-lenpasswd-file` flag.
```
user:password
admin:qwerty123
bob:wbxo28skc
```

Yes, the passwords in the file are stored in unencrypted form.

Now when you try to create a paste via a web browser or API you will be prompted for authorization.
