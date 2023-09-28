use tonic::transport::Channel;
use tonic::Request;

pub mod user {
    tonic::include_proto!("user");
}

use user::user_service_client::UserServiceClient;
use user::GetUserDetailsRequest;
use user::UpdateUserDetailsRequest;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let channel = Channel::from_static("http://localhost:50051")
        .connect()
        .await?;
    let mut client = UserServiceClient::new(channel);
    let request = Request::new(GetUserDetailsRequest {
        username: "user2".into(),
    });
    let response = client.get_user_details(request).await?;
    let response = response.into_inner();
    println!("RESPONSE={:?}", response);
    // now use the response to create an update request
    let mut user = response.user_data.unwrap();
    user.last_name = "updated from Rust".into();
    let request = Request::new(UpdateUserDetailsRequest {
        user_data: Some(user),
    });
    let response = client.update_user_details(request).await?;
    let response = response.into_inner();
    println!("RESPONSE={:?}", response);
    Ok(())
}
