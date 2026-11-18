# Layout actuel

vertical_padding
title (responsive)
track info
timestatus
spectrum (only in mainMode, height responsive)
seekbar
settingsStacked (only in mainMode, height responsive)
body
hepbar
bottomBar
transient
vertical_padding

# Nouveau layout

vertical_padding
title (responsive)
track info
timestatus
spectrum (only in mainMode, height responsive)
seekbar
settingsStacked (only in mainMode, height responsive)
gap (new section with empty line, responsive)
body
transient
hepbar
vertical_padding

(bottom bar has been removed)


# Répartition en hauteur

No need to name layouts (compact, minimal etc..). Only two layout will remain: oneColumn and twoColumn (see below)

Each section have it's height defined, they body takes remaining space

layoutHeight = min(terminal height - 2*vertical_padding, 20 + vis_rows)
if (expanded) layoutHeight = terminal height - 2*vertical_padding

le layoutHeight commun quelque soit la page, ça ne changera jamais

titleHeight = layoutHeight < 25 ? 0 : 1

spectrumHeight = 0
if (queueView)
  if layoutHeight < 25 spectrumHeight = 0
  if layoutHeight < 30 spectrumHeight = 4
  else spectrumHeight = config.vis_rows (default 7)

helbarHeight = helpBarVisible ? 1 : 0

gapHeight = layoutHeight < 30 ? 0 : 1

trackInfoHeight, TimeStatusHeight, seekBarHeight, TransientHeight = toujours à 1

bodyHeight = layoutHeight - titleHeight - trackInfoHeight - TimeStatusHeight - seekBarHeight - SettingsStackedHeight - gapHeight - TransientHeight - helpBarHeight

-> le body prend toute la place verticale disponible, meme si il ne contient que quelques lignes. Donc le transient et helpbar seront toujours collé en bas

# Répartition horizontale du body en **mode Queue**

layoutWidth = terminalWidth - 2xhorizontal_padding

## Settings

Settings = SRC, VOL, EQ, SHF, RPT, SPD, NET

SettingsStacked = rows below seekbar to display the settings
SettingsPane = column in side bar

## if layoutWidth >= 80 && settingsPane activated

-> Two column layout
SettingsColoumnWidth = 30
GridGapWidth = 5
bodyWidth = layoutWidth - SettingsColoumnWidth - GridGapWidth
SettingsStackedHeight = 0

## else if layoutWidth > 90

-> one column layout but with large width. All settings fit in one rows
SettingsStackedHeight = 1
bodyWidth = layoutWidth

## else

-> one column layout but with small width. Settings need two rows
SettingsStackedHeight = 2
bodyWidth = layoutWidth