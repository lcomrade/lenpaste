# Themes for WEB interface
## Introduction
Lenpaste has two built-in themes: dart and light.
You can switch between the themes in the settings.

You can also add your own themes by placing them in the `/data/themes/` directory in Docker.
If you work with bare metal, specify the path to the themes folder using the `-ui-themes-dir` flag.



## Specife default server theme
You can also specify the default theme with the `LENPASTE_UI_DEFAULT_THEME` variable in Docker.
Or the `-ui-default-theme` flag for bare metal.

The value of this parameter must match the name of the theme file without an extension.
For example: `dark`, `light`, `my_theme`.



## Create custom theme
Use one of the themes from the `./internal/web/data/theme/` directory as a base.
The extension `.theme` in the file name is mandatory.

If you miss a theme parameter, it will be replaced by the default one.
