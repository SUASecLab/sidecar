from setuptools import setup

setup(
    name="sidecar_lib",
    version="0.0.1",
    description="sidecar lib for the SUASecLab sidecar",
    url="https://github.com/SUASecLab/sidecar",
    author="Tobias Tefke",
    author_email="t.tefke@stud.fh-sm.de",
    license="GPLv3.0", 
    packages=["sidecar_lib"],
    install_requires=["web.py",
                 "pyjwt",
                 "pymongo",
                 ],
    classifiers=[],
)
