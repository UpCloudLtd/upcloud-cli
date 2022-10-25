import os

def withBase(name, base):
    if base:
        return f'{base} {name}'
    return name

def parse_name(name):
    return name.replace('_', ' ').replace('.md', '').replace('index', '')

def generateNav(path, base=''):
    pages = []

    for i in sorted(os.scandir(path), key=lambda i: parse_name(i.name)):
        if i.is_file():
            name = parse_name(i.name)
            pages.append({withBase(name, base): i.path})
        if i.is_dir():
            name = i.name.replace('_', ' ')
            pages.append({withBase(name, base): generateNav(i.path, base=name.replace('upctl ', ''))})

    return pages

if __name__ == '__main__':
    os.chdir('docs/')
    nav = generateNav('commands_reference/')
    os.chdir('..')

    with open("mkdocs.base.yaml") as f:
        config = f.read().replace('Commands reference: []', f'Commands reference: {nav}')
    
    with open("mkdocs.yaml", "w") as f:
        f.write(config)
