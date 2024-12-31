### GOLANG-JWT-AUTHENTICATION

- This a RESTful API that uses JWT for authentication.
- It is built with Golang and the Gin framework.
- It 2 main routes which are:

  - `/api/auth/signup` - for registering a user
  - `/api/auth/login` - for logging in a user

- The API uses refresh tokens to generate new access tokens when the current access token expires.
- Maximum 5 valid sessions are allowed per user.
- The API uses MongoDB as the database.

<p>Try this API service in postman</p>

[<img src="https://run.pstmn.io/button.svg" alt="Run In Postman" style="width: 128px; height: 32px;">](https://app.getpostman.com/run-collection/21279139-42d703df-2115-43fe-adbb-2f6d6deb3c89?action=collection%2Ffork&source=rip_markdown&collection-url=entityId%3D21279139-42d703df-2115-43fe-adbb-2f6d6deb3c89%26entityType%3Dcollection%26workspaceId%3D43ffaeaf-502b-40d6-a91f-8727aa39a837)
