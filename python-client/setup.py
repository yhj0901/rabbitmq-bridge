from setuptools import setup, find_packages

setup(
    name="rabbitmq-bridge-client-yhj",
    version="0.1.0",
    packages=find_packages(),
    install_requires=[
        "pika>=1.3.0",
    ],
    author="heejune Yang",
    author_email="yangheejune0901@gmail.com",
    description="RabbitMQ client library for messaging between Go and Python",
    long_description=open("README.md").read(),
    long_description_content_type="text/markdown",
    url="https://github.com/yhj0901/rabbitmq-bridge",
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
    ],
    python_requires=">=3.7",
) 