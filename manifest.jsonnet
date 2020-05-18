local manifest = import 'core/manifest.libsonnet';
local icons() = {
  [size]: 'icon.png'
  for size in ['16', '48', '128']
};

local json = manifest.new(
  name='Go Search Extension',
  version='0.2',
  keyword='go',
  description='Search Golang std docs and third packages in your address bar instantly!',
)
             .addIcons(icons())
             .addBackgroundScripts([
  'search/docs.js',
  'search/package.js',
  'index/godocs.js',
  'index/packages.js',
  'main.js',
]);

if std.extVar('browser') == 'firefox' then
  json
  + {
    browser_specific_settings: {
      gecko: {
        id: '{bb3394f3-9a10-4b30-ae77-6cc6fd51de99}',
      },
    },
  }
else
  json
