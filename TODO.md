### ✅ Refactoring list
### Goals
>- Implement Hexagonal/Onion architecture
>- Pass `https://goreportcard.com/` validation
>- Make code more robust

## 🔴 Refactoring Tasks To-Do High Priority (1)
- [ ] Make application runnable even if there is no other running services (mongo/rabbit/etc...)
- [x] Implemnet Repository instead of explicit usage of mongodb
- [ ] Remove config usage from all modules and move it to the `main`
- [ ] Move code under appropropriate `internal` and `pkg` dirs
- [ ] Remove `Poor Packaging` such as directory `model`,`util`...
- [ ] Move middleware creation from `main` to appropriate modules
- [ ] Make connector interfaces `Ports`
- [ ] Reduce Coupling in code
- [ ] Reduce Coupling in data structures
- [ ] Minimize amount of dynamically allocated object `&` or `new()`
- [ ] Generate endpoints automatically via `oapi-codegen`, instead of specifying them manually
- [ ] Instead of handler struct, use CQRS for better code separation
- [ ] Improve registration to handle both regular auth and OAuth without using a boolean to skip password check

## 🔴 Implementation Tasks (1.5) 
- [ ] Implement mechanism to send tasks only to specific workers
- [x] Change Rabbitmq payload type from raw binary to Protobuffs - *Changed: JSON*
- [ ] Define strict rules for workers discovery
- [ ] Implement mechanism to send tasks only to specific workers
- [ ] Increase test coverage
- [ ] Add [golangci-lint](https://golangci-lint.run/)


## 🟠 To-Do Medium Priority (2)
- [ ] Measure passing by value vs ref performance
- [ ] Add ability to be logged in one account from different devices
- [x] Change `GetBasics` to `GetUserInfo`, it seems more obvious
- [ ] Change `fullname` to `nickname`, because service do not require real name

## 🔵 Backlog (0)
- [c] Both `GetUserByEmail` and `GetUserById` may be redundand, consider refactoring - *Canceled: both useful*
