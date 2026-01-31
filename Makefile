.PHONY: mock
mock:
	@mockgen -source=E:\Project\my_project\WeBook\webook\internal\service\user.go -package=svcmocks -destination=E:\Project\my_project\WeBook\webook\internal\service\mocks\user.mock.go
	@mockgen -source=E:\Project\my_project\WeBook\webook\internal\service\code.go -package=svcmocks -destination=E:\Project\my_project\WeBook\webook\internal\service\mocks\code.mock.go
	@mockgen -source=E:\Project\my_project\WeBook\webook\internal\repository\user.go -package=repomocks -destination=E:\Project\my_project\WeBook\webook\internal\repository\mocks\user.mock.go
	@mockgen -source=E:\Project\my_project\WeBook\webook\internal\repository\code.go -package=repomocks -destination=E:\Project\my_project\WeBook\webook\internal\repository\mocks\code.mock.go
	@mockgen -source=E:\Project\my_project\WeBook\webook\internal\repository\dao\user.go -package=daomocks -destination=E:\Project\my_project\WeBook\webook\internal\repository\dao\mocks\user.mock.go
	@mockgen -source=E:\Project\my_project\WeBook\webook\internal\repository\cache\user.go -package=cachemocks -destination=E:\Project\my_project\WeBook\webook\internal\repository\cache\mocks\user.mock.go
	@mockgen -package=redismocks -destination=E:\Project\my_project\WeBook\webook\internal\repository\cache\redismocks\cmd.mock.go github.com/redis/go-redis/v9 Cmdable
	@go mod tidy