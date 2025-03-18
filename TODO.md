### ✅ Refactoring list
### Goals
>- Implement Hexagonal/Onion architecture
>- Pass `https://goreportcard.com/` validation
>- Make code more robust

## 🔴 Refactoring Tasks To-Do High Priority (1)
- [ ] Remove config usage from all modules and move it to the `main`
- [ ] Move code under appropropriate `internal` and `pkg` dirs
- [ ] Remove `Poor Packaging` such as directory `model`,`util`...
- [ ] Move middleware creation from `main` to appropriate modules
- [ ] Make connector interfaces `Ports`
- [ ] Reduce Coupling in code
- [ ] Reduce Coupling in data structures
- [ ] Minimize amount of dynamically allocated object `&` or `new()`

## 🔴 Implementation Tasks (1.5) 
- [ ] Implement mechanism to send tasks only to specific workers
- [ ] Change Rabbitmq payload type from raw binary to Protobuffs
- [ ] Define strict rules for workers discovery
- [ ] Implement mechanism to send tasks only to specific workers
- [ ] Increase test coverage

## 🟠 To-Do Medium Priority (2)
- [ ] Measure passing by value vs ref performance
- [ ] Add ability to be logged in one account from different devices