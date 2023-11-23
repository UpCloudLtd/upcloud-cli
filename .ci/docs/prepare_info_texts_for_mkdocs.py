'''Prepare info texts for mkdocs

This script overwrites code fences in the .md files copied into docs/examples so that only language and title are left in-place, because mkdocs (and the used theme) has limited support for fenced code-block info texts.
'''
import fileinput
import os
import re
import sys

def get_markdown_files(path):
    return (f.path for f in os.scandir(path) if f.name.endswith('.md'))

def simplify_info_texts(path):
    for line in fileinput.input(path, inplace=True):
        if line.startswith('```'):
            start = re.search(r'^```[a-zA-Z]*', line)
            title = re.search(r'title=".*"', line)

            old = line.rstrip('\n')
            new = f'{start.group(0)} {title.group(0)}' if title else f'{start.group(0)}'

            if new != old:
                sys.stderr.write(f'{fileinput.filename()}:{fileinput.filelineno()}:\n- {old}\n+ {new}\n')
            print(new)
        else:
            print(line, end='')

if __name__ == '__main__':
    files = get_markdown_files('docs/examples/')
    for f in files:
        simplify_info_texts(f)
