package configs

//Const config system
const (
	// /neo4jURL = "Bolt://neo4j:tlis2016@tlu.cloudapp.net:7687"
	URLDB   = "http://neo4j:madawg00@tlu.cloudapp.net:7474/db/data/"
	APIPort = "8080"
)

// Security const
const (
	JWTSecretKey = "socialnetworkTLU"
)

//ObjectType const
const (
	PostObject    = 1
	UserObject    = 2
	CommentObject = 3
	LikeObject    = 4
	PageObject    = 5
	GroupObject   = 6
)

// Const privacy
const (
	Public           = 1
	ShareToFollowers = 2
	Private          = 3
)

// ErrorCode Table
const (
	EcSuccess                           = 1   //	Success
	EcNoExistObject                     = 2   //	No exist this object.
	EcParam                             = 100 //	Invalid parameter
	EcParamMissingField                 = 101 //	Missing a few fields.
	EcParamUserID                       = 110 //	Invalid user id
	EcParamUserField                    = 111 //	Invalid user info field
	EcParamEmail                        = 113 //	Invalid email
	EcParamFieldList                    = 115 //	Invalid field list
	EcParamPhotoID                      = 121 //	Invalid photo id
	EcParamTitle                        = 142 //	Invalid title
	EcParamAccessToken                  = 190 //	Invalid OAuth 2.0 Access Token
	EcPermission                        = 200 //	Permissions error
	EcPermissionUser                    = 210 //	User not visible
	EcPermissionPhoto                   = 221 //	Photo not visible
	EcPermissionMessage                 = 230 //	Permissions disallow message to user
	EcEdit                              = 300 //	Edit failure
	EcEditUserData                      = 310 //	User data edit failure
	EcUsersCreateInvalidEmail           = 370 //	The email address you provided is not a valid email address
	EcUsersCreateExistingEmail          = 371 //	The email address you provided belongs to an existing account
	EcUsersCreateBirthday               = 372 //	The birthday provided is not valid
	EcUsersCreatePassword               = 373 //	The password provided is too short or weak
	EcUsersRegisterInvalidCredential    = 374 //	The login credential you provided is invalid.
	EcUsersRegisterConfFailure          = 375 //	Failed to send confirmation message to the specified login credential.
	EcUsersRegisterExisting             = 376 //	The login credential you provided belongs to an existing account
	EcUsersRegisterDefaultError         = 377 //	Sorry, we were unable to process your registration.
	EcUsersRegisterPasswordBlank        = 378 //	Your password cannot be blank. Please try another.
	EcUsersRegisterPasswordInvalidChars = 379 //	Your password contains invalid characters. Please try another.
	EcUsersRegisterPasswordShort        = 380 //	Your password must be at least 6 characters long. Please try another.
	EcUsersRegisterPasswordWeak         = 381 //	Your password should be more secure. Please try another.
	EcUsersRegisterUsernameError        = 382 //	Please enter a valid username.
	EcUsersRegisterMissingInput         = 383 //	You must fill in all of the fields.
	EcUsersRegisterIncompleteBday       = 384 //	You must indicate your full birthday to register.
	EcUsersRegisterInvalidEmail         = 385 //	Please enter a valid email address.
	EcUsersRegisterEmailDisabled        = 386 //	The email address you entered has been disabled. Please contact disabled@facebook.com with any questions.
	EcUsersRegisterAddUserFailed        = 387 //	There was an error with your registration. Please try registering again.
	EcUsersRegisterNoGender             = 388 //	Please select either Male or Female.
	EcAuthEmail                         = 400 //	Invalid email address
	EcAuthLogin                         = 401 //	Invalid username or password
	EcAuthMissingToken                  = 404 //	Missing token.
	EcAuthInvalidToken                  = 405 //	Invalid token.
	EcAuthNoExistToken                  = 406 //	No exist token.
	EcAuthCheckToken                    = 407 //	Error in checking token.
	EcAuthGenerateToken                 = 408 //	Error in generate token.
	EcAuthNoExistUser                   = 409 //	No exist user.
	EcAuthNoExistFacebook               = 410 //	No exist account with this facebook.
	EcAuthInvalidFacebookToken          = 411 //	Error in checking token.
	EcAuthWrongPassword                 = 412 //	Error in login: Wrong password.
	EcAuthNoExistEmail                  = 413 //
	EcAuthWrongRecoveryCode             = 414 //	Error in recover password: Wrong recovery code.
	EcMesgNoBody                        = 501 //	Missing message body
)

// TypePost const
const (
	Post           = 0
	PostStatus     = 1
	PostPhoto      = 2
	PostLink       = 3
	PostGroup      = 4
	PostSharePost  = 5
	PostSharePage  = 6
	PostShareGroup = 7
)

//FCMToken struct
const (
	FCMToken = "AAAAuET9LvY:APA91bEYl-fIkcY0w7b6umgBHD4yrZnG_v9I2iY1K3EnjUfSrYvlFYIG5vrmP8wFCH8ZMZ-Kx6U6u3XIsw-AIGehs-msWXtlzOq8R_50qAiqcsrJv9WQluALvjWPqSIAPrVS2RKZ4H6V"
)
