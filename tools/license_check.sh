#!/bin/sh
set -e

eval failed=false

LICENSE_GO_SED='1,18p'
LICENSE_GO=$(cat <<EOF
// Copyright (C) 2021-$(date +%Y) Leonid Maslakov.

// This file is part of Lenpaste.

// Lenpaste is free software: you can redistribute it
// and/or modify it under the terms of the
// GNU Affero Public License as published by the
// Free Software Foundation, either version 3 of the License,
// or (at your option) any later version.

// Lenpaste is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
// or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero Public License for more details.

// You should have received a copy of the GNU Affero Public License along with Lenpaste.
// If not, see <https://www.gnu.org/licenses/>.
EOF
)

while read -r file; do
    if [ "$(sed -n "${LICENSE_GO_SED}" "${file}")" != "${LICENSE_GO}" ]; then
        echo "license check failed: $file"
    fi
done <<EOT
$(find ./ -type f -name '*.go' ! -name '*_test.go')
EOT

while read -r file; do
    if [ "$(sed -n "${LICENSE_GO_SED}" "${file}")" != "${LICENSE_GO}" ]; then
        failed=true
        echo "license check failed: $file"
    fi
done <<EOT
$(find ./ -type f -name '*.js')
EOT



LICENSE_TMPL_SED='1,20p'
LICENSE_TMPL=$(cat <<EOF
{{/*
   Copyright (C) 2021-$(date +%Y) Leonid Maslakov.

   This file is part of Lenpaste.

   Lenpaste is free software: you can redistribute it
   and/or modify it under the terms of the
   GNU Affero Public License as published by the
   Free Software Foundation, either version 3 of the License,
   or (at your option) any later version.

   Lenpaste is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
   or FITNESS FOR A PARTICULAR PURPOSE.
   See the GNU Affero Public License for more details.

   You should have received a copy of the GNU Affero Public License along with Lenpaste.
   If not, see <https://www.gnu.org/licenses/>.
*/}}
EOF
)

while read -r file; do
    if [ "$file" = "./internal/handler/data/license.tmpl" ]; then
        continue
    fi

    if [ "$(sed -n "${LICENSE_TMPL_SED}" "${file}")" != "${LICENSE_TMPL}" ]; then
        failed=true
        echo "license check failed: $file"
    fi
done <<EOT
$(find ./ -type f -name '*.tmpl')
EOT



LICENSE_CSS_SED='1,20p'
LICENSE_CSS=$(cat <<EOF
/*
   Copyright (C) 2021-$(date +%Y) Leonid Maslakov.

   This file is part of Lenpaste.

   Lenpaste is free software: you can redistribute it
   and/or modify it under the terms of the
   GNU Affero Public License as published by the
   Free Software Foundation, either version 3 of the License,
   or (at your option) any later version.

   Lenpaste is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
   or FITNESS FOR A PARTICULAR PURPOSE.
   See the GNU Affero Public License for more details.

   You should have received a copy of the GNU Affero Public License along with Lenpaste.
   If not, see <https://www.gnu.org/licenses/>.
*/
EOF
)

while read -r file; do
    if [ "$(sed -n "${LICENSE_CSS_SED}" "${file}")" != "${LICENSE_CSS}" ]; then
        failed=true
        echo "license check failed: $file"
    fi
done <<EOT
$(find ./ -type f -name '*.css')
EOT



if [ $failed = true ]; then
    exit 1
fi
