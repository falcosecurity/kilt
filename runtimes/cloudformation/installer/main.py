import click
import boto3
import threading
import logging
import sys
from os import path
from jinja2 import Environment, FileSystemLoader
from botocore.exceptions import ClientError


bundle_dir = path.dirname(__file__)
# If inside pyinstaller we use temp dir as bundle dir
if getattr(sys, 'frozen', False) and hasattr(sys, '_MEIPASS'):
    bundle_dir = getattr(sys, '_MEIPASS', path.abspath(path.dirname(__file__)))

KILT_CFN = path.join(bundle_dir, 'kilt.yaml')
KILT_ZIP = path.join(bundle_dir, 'kilt.zip')

assert path.exists(KILT_CFN), 'Could not find cloudformation jinja template'
assert path.exists(KILT_ZIP), 'Could not find kilt.zip - did you build it?'


class CallbackProgress:
    def __init__(self, filename, label):
        self._size = path.getsize(filename)
        self._lock = threading.Lock()
        self._bar = click.progressbar(label=label, length=self._size)
        self._bar.render_progress()

    def __call__(self, bytes_amount):
        with self._lock:
            self._bar.update(bytes_amount)

    def done(self):
        self._bar.render_finish()


@click.command('kilt-cfn-installer')
@click.argument('macro_name')
@click.argument('path_to_kilt_definition', type=click.Path(exists=True))
@click.option('--region', '-r', help="override aws region")
@click.option('--opt-in/--opt-out', is_flag=True, help="Use opt-in mechanism instead of the default opt out")
@click.option('--kilt-zip-name', default="kilt.zip", help="[Optional] Name of the file for the lambda code")
@click.option('--kms-secret', default="",
              help="[Optional] ARN of the secret containing credentials for the image repository")
@click.option('--kilt-zip', default=KILT_ZIP, help='Deploy custom lambda instead of bundled lambda')
@click.option('--kilt-template', default=KILT_CFN, help='Use custom CFN template instead of bundled one')
def main(macro_name, path_to_kilt_definition, region, opt_in, kilt_zip_name, kms_secret, kilt_zip, kilt_template):
    click.echo("Getting AWS account and region...", nl=False)
    aws_account = boto3.client('sts').get_caller_identity().get('Account')
    aws_region = boto3.session.Session().region_name
    if region is not None:
        aws_region = region
    click.echo(click.style(f'{aws_account} in {aws_region}', fg='green'))
    click.echo("Getting S3 bucket...", nl=False)
    s3_bucket = get_s3_bucket(aws_account, aws_region)
    click.echo(click.style(s3_bucket.name, fg='green'))
    pb = CallbackProgress(kilt_zip, "Uploading Macro")
    s3_bucket.upload_file(kilt_zip, kilt_zip_name, Callback=pb)
    pb.done()
    pb = CallbackProgress(path_to_kilt_definition, f"Uploading Kilt - {macro_name}")
    s3_bucket.upload_file(path_to_kilt_definition, f'{macro_name}.kilt.cfg', Callback=pb)
    pb.done()
    env = Environment(
        loader=FileSystemLoader(searchpath=path.dirname(kilt_template))
    )
    template = env.get_template(path.basename(kilt_template))
    output_text = template.render(
        macro_name=macro_name,
        bucket_name=s3_bucket.name,
        kilt_zip_name=kilt_zip_name,
        kilt_opt_in=opt_in,
        kilt_kms_secret=kms_secret
    )
    cf = boto3.resource('cloudformation', region_name=aws_region)
    stack_name = f'KiltMacro{macro_name}'
    click.echo(f"Creating stack {stack_name}...", nl=False)
    stack = cf.create_stack(
        StackName=stack_name,
        TemplateBody=output_text,
        Capabilities=[
            'CAPABILITY_IAM'
        ]
    )
    click.echo(click.style("SUBMITTED", fg="yellow"))
    click.echo(f"Check your cloudwatch console for stack {stack_name}")


def get_s3_bucket(aws_account, aws_region):
    try:
        s3_bucket_name = f'kilt-{aws_account}-{aws_region}'
        s3 = boto3.resource(service_name='s3', region_name=aws_region)
        location = {'LocationConstraint': aws_region}
        s3.create_bucket(
            ACL='private',
            Bucket=s3_bucket_name,
            CreateBucketConfiguration=location
        )

        bucket = s3.Bucket(s3_bucket_name)
        bucket.wait_until_exists()
    except ClientError as e:
        logging.error(e)
        sys.exit(1)
    return bucket


if __name__ == '__main__':
    main()
