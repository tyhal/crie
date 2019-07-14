# Crie enforcement

[![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/tyhal/crie.svg)](https://hub.docker.com/r/tyhal/crie)

Crie is an effective way to format and lint code for a variety of languages

-   No dependencies ( except docker which all build machines have )
-   Small docker image reduces the wait time and disk usage ( &lt; 500MB)
-   Extendable for more languages
-   Portable for installation as the run command will download it

```bash
    sudo ./script/crie install
```

### Crie Status

Required

| name              | chk                | fmt                |
| ----------------- | ------------------ | ------------------ |
| **Language**      |                    |                    |
| Shell             | :white_check_mark: | :white_check_mark: |
| Cpp               | :white_check_mark: | :white_check_mark: |
| CppCuda           | :x:                | :x:                |
| Go                | :white_check_mark: | :white_check_mark: |
| Groovy            | :x:                | :x:                |
| JavaScript +JSX   | :white_check_mark: | :white_check_mark: |
| Python            | :white_check_mark: | :white_check_mark: |
| **Config**        |                    |                    |
| CSS + LESS + SASS | :x:                | :x:                |
| YML               | :white_check_mark: | :x:                |
| JSON              | :white_check_mark: | :white_check_mark: |
| **Build**         |                    |                    |
| Docker            | :white_check_mark: | :x:                |
| CMake             | :white_check_mark: | :x:                |
| Compose           | :white_check_mark: | :x:                |
| **Configuration** |                    |                    |
| Ansible           | :x:                |                    |
| Terraform         |                    | :white_check_mark: |
| **Documentation** |                    |                    |
| Markdown          | :white_check_mark: | :white_check_mark: |
| Doxygen           | :white_check_mark: | :x:                |
| **Other**         |                    |                    |
| Project Structure | :white_check_mark: |                    |

Extra

| name    | chk | fmt |
| ------- | --- | --- |
| Haskell | :x: | :x: |
| Java    | :x: | :x: |
| Julia   | :x: | :x: |
| Rust    | :x: | :x: |
| Kotlin  | :x: | :x: |
