# -*- mode: python ; coding: utf-8 -*-
import os

block_cipher = None


a = Analysis(['main.py'],
             pathex=[os.getcwd()],
             binaries=[],
             datas=[
                 ('../kilt.zip', '.'),
                 ('./kilt.yaml', '.'),

             ],
             hiddenimports=[],
             hookspath=[],
             runtime_hooks=[],
             excludes=[],
             win_no_prefer_redirects=False,
             win_private_assemblies=False,
             cipher=block_cipher,
             noarchive=False)
pyz = PYZ(a.pure, a.zipped_data,
             cipher=block_cipher)
exe = EXE(pyz,
          a.scripts,
          a.binaries,
          a.zipfiles,
          a.datas,
          [],
          name='kilt-installer',
          debug=False,
          bootloader_ignore_signals=False,
          strip=True,
          upx=False,
          upx_exclude=[],
          runtime_tmpdir=None,
          console=True )
