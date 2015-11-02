# projectmain

Rename this repo to the actual project name (not May1601 - something cool).

### Design

**Note**: the plugin system is *heavily* inspired by and based on the
excellent plugin architecture design of
[packer](https://github.com/mitchellh/packer).


#### Project Structure

`projectmain`: core interfaces

`projectmain/rpc`: core interface implementations to be executed over RPC

`projectmain/plugin`: sets up plugins by serving a particular interface
