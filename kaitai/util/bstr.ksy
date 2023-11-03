# SPDX-License-Identifier: MIT

meta:
  id: bstr
  endian: le

seq:
  - id: len
    type: u4
  - id: str
    type: str
    encoding: UTF-16LE
    size: len
