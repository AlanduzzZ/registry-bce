import os
from jinja2 import Environment, FileSystemLoader, StrictUndefined
from utils.migration import read_conf

revision = '2.2.0'
down_revisions = ['2.1.0']

def migrate(input_cfg, output_cfg):
    current_dir = os.path.dirname(__file__)
    tpl = Environment(
        loader=FileSystemLoader(current_dir),
        undefined=StrictUndefined,
        trim_blocks=True,
        lstrip_blocks=True
        ).get_template('harbor.yml.jinja')

    config_dict = read_conf(input_cfg)

    with open(output_cfg, 'w') as f:
        f.write(tpl.render(**config_dict))
