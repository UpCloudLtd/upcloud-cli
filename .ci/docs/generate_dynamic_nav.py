import os

def withBase(name, base):
    if base:
        return f'{base} {name}'
    return name

def parse_name(name):
    return name.replace('_', ' ').replace('.md', '').replace('index', '')

def generateCommandsReferenceNav(path, base=''):
    pages = []

    for i in sorted(os.scandir(path), key=lambda i: parse_name(i.name)):
        if i.is_file():
            name = parse_name(i.name)
            pages.append({withBase(name, base): i.path})
        if i.is_dir():
            name = i.name.replace('_', ' ')
            pages.append({withBase(name, base): generateCommandsReferenceNav(i.path, base=name.replace('upctl ', ''))})

    return pages

def generateExamplesNav(path):
    return [i.path for i in sorted(os.scandir(path), key=lambda i: i.name) if i.name != 'README.md']


if __name__ == '__main__':
    navs = dict()

    os.chdir('docs/')
    navs["Commands reference"] = generateCommandsReferenceNav('commands_reference/')
    navs["Examples"] = generateExamplesNav('examples/')
    os.chdir('..')

    with open("mkdocs.base.yaml") as f:
        config = f.read()

        for key, nav in navs.items():
            config = config.replace(f'{key}: []', f'{key}: {nav}')

    with open("mkdocs.yaml", "w") as f:
        f.write(config)
