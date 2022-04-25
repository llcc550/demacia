# v0.0.1 (2021.12.13)

### Added

- 完成websocket服务。 demo为student.api的导入服务
- 接入分布式事务dtm 。 demo为demo为auth.api分别调用user.rpc和member.rpc
- 完成demo。涉及服务student.api、class.api、class.rpc、member.rpc、user.rpc、organization.rpc、auth.api
- 统一api错误类型

# v0.0.2 (2021.12.14)

### Changed

- websocket所需的消息队列，由beanstalkd切换为kafka

# v0.0.3 (2022.01.03)

### Upgrade

- go-zero升级到1.2.5

# v0.0.4 (2022.02.08)

### Upgrade

- 添加gitlab-runner
- 拆分错误信息
- 合并member服务
- 合并card服务