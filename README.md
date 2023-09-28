## Rust and Go clients to Go server

This is an example of how a Rust client written using [tonic] might make us lose data when communicating with a Go server.

The server and client are exchanging a simple message:

```protobuf
message UserData {
  string username = 1;
  string first_name = 2;
  string last_name = 3;
  string email = 4; // New field added later
}
```

The server has been updated with the new field, while the clients are using the old version of the message.
For reference, this is the diff between the two versions:

```diff
diff -Nurp go-service/proto/user/user.proto rust-client/proto/user/user.proto
--- go-service/proto/user/user.proto    2023-09-28 15:50:46
+++ rust-client/proto/user/user.proto   2023-09-28 15:33:05
@@ -8,7 +8,7 @@ message UserData {
   string username = 1;
   string first_name = 2;
   string last_name = 3;
-  string email = 4; // New field added only in Go service
+  // string email = 4; // New field added only in Go service
 }
 ```

 The Go server keeps an array in memory with fake data, just a few users:
 ```
 ❯ ./go-service/bin/service
Starting server...
Username:  user1
First name:  John
Last name:  Doe
Email:  user1@domain.com
-----
Username:  user2
First name:  Jane
Last name:  Doe
Email:  user2@domain.com
-----
2023/09/28 16:26:09 Server listening on port 50051
 ```

Now we use the Go client to find the `user2` user and update its Last name to the literal string `"updated"`:

```
❯ ./go-client/bin/client
2023/09/28 16:27:08 Response received: user_data:{username:"user2" first_name:"Jane" last_name:"Doe" 4:"user2@domain.com"}
2023/09/28 16:27:08 Response received: user_data:{username:"user2" first_name:"Jane" last_name:"updated" 4:"user2@domain.com"}
```

This is the server log for that transaction:

```
2023/09/28 16:27:08 Received request: username:"user2"
Found user:  user2
2023/09/28 16:27:08 Received request: user_data:{username:"user2" first_name:"Jane" last_name:"updated" email:"user2@domain.com"}
Username:  user1
First name:  John
Last name:  Doe
Email:  user1@domain.com
-----
Username:  user2
First name:  Jane
Last name:  updated
Email:  user2@domain.com
-----
```

Note that `user2` has an updated last name now.

Now we run the Rust client to perform the same opration,
updating the last name of `user2` to the literal string `"updated from Rust"`:

```
❯ ./rust-client/target/debug/rust-client
RESPONSE=GetUserDetailsResponse { user_data: Some(UserData { username: "user2", first_name: "Jane", last_name: "updated" }) }
RESPONSE=UpdateUserDetailsResponse { user_data: Some(UserData { username: "user2", first_name: "Jane", last_name: "updated from Rust" }) }
```

This is the server log for that transaction:

```
2023/09/28 16:29:16 Received request: username:"user2"
Found user:  user2
2023/09/28 16:29:16 Received request: user_data:{username:"user2" first_name:"Jane" last_name:"updated from Rust"}
Username:  user1
First name:  John
Last name:  Doe
Email:  user1@domain.com
-----
Username:  user2
First name:  Jane
Last name:  updated from Rust
Email:
-----
```

Note that we've lost the email field for `user2`!

The problem is that [tonic] uses [prost] under the covers and [prost] doesn't support unknown fields.
See [this issue](https://github.com/tokio-rs/prost/issues/2) for more details, which is linked form the
results of running the [protobuf conformance test suite against prost](https://github.com/tokio-rs/prost/blob/master/conformance/failing_tests.txt).

[tonic]: https://github.com/hyperium/tonic
[prost]: https://github.com/tokio-rs/prost
