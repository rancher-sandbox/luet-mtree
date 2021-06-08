# luet-mtree

luet-mtree is a plugin for [luet](https://luet-lab.github.io/docs/) to create and check [mtree](https://www.freebsd.org/cgi/man.cgi?mtree(8)) validation files.
Currently is in use by [cOS-toolkit](https://github.com/rancher-sandbox/cOS-toolkit) during packages build to generate the checksums and to validate those packages during system upgrades



## Commands

### check

Checks a file or dir against a validation file.
This command supports parsing a validation file against a tar archive uncompressed or compressed with zstd/gzip instead of a file/dir

```bash
Usage:
  luet-mtree check [file or dir] [validation file] [flags]

Flags:
  -x, --exclude strings   Exclude paths from check. Checks against the path prefix, so 'oem/' will cover both 'oem/' and 'oem/features/' paths.
  -f, --format string     Format for output. Choices are bsd, path and json. (default "bsd")
  -h, --help              help for check
```

### generate

Generate a checksum file for the file or dir provided
This command supports generating a validation file from a tar archive uncompressed or compressed with zstd/gzip instead of a file/dir

```bash
Usage:
  luet-mtree generate [file or dir] [flags]

Flags:
  -h, --help               help for generate
  -k, --keywords strings   Keywords to use to generate the tree (sha256 will automatically be added)
  -o, --output string      Name for output file, otherwise it defaults to stdout
```

### event parsing

When used as a plugin to luet (by calling luet with `--plugin luet-mtree`, see [plugin docs](https://luet-lab.github.io/docs/docs/concepts/plugins-and-extensions/)) luet emits events based on the actions being performed.
We take those events and payloads and execute an action depending on their contents.

You can manually test those events by calling luet-mtree with no subcommands and 2 params. The first being the event emitted (see events emmited by luet [here](https://github.com/mudler/luet/blob/master/pkg/bus/events.go)) and the second a json payload, the contents depend on the type of event.

```bash
luet-mtree 'image.post.unpack' '{"data": {"Image": "quay.io/costoolkit/releases-opensuse:systemd-boot-live-26", "Dest": "/tmp/upgrade"}}'
```


## License

Copyright (c) 2021 [SUSE, LLC](http://suse.com)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.