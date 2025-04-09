package scanner

// Compacted form of the code developed by Daniel Alder
// https://github.com/daald/dpkg-licenses

var dpkg_cmd = `
set -e

machine_readable () {
  package="$1"
  copyrightfile=
  if [ -f "/usr/share/doc/$package/copyright" ]; then
    copyrightfile="/usr/share/doc/$package/copyright"
  elif [ -f "/usr/share/doc/${package%:*}/copyright" ]; then
    copyrightfile="/usr/share/doc/${package%:*}/copyright"
  else
    exit 0  # no copyright file found
  fi

  format=$(awk '/^Format:/{print}/^$/{exit}' "$copyrightfile")
  [ -n "$format" ] || exit 0

  case "$format" in
    *'//www.debian.org/doc/packaging-manuals/copyright-format/1.0'*)
      result=$(grep '^License:' "$copyrightfile" | cut -d':' -f2-)
      ;;
    *'//dep.debian.net/deps/dep5'*)
      result=$(grep '^License:' "$copyrightfile" | cut -d':' -f2-)
      ;;
    *'//anonscm.debian.org/viewvc/dep/web/deps/dep5.mdwn?'*)
      result=$(grep '^License:' "$copyrightfile" | cut -d':' -f2-)
      ;;
    *'//svn.debian.org/wsvn/dep/web/deps/dep5.mdwn?'*)
      result=$(grep '^License:' "$copyrightfile" | cut -d':' -f2-)
      ;;
    *'//anonscm.debian.org/loggerhead/dep/dep5/trunk/annotate/179/dep5/copyright-format.xml'*)
      result=$(grep '^License:' "$copyrightfile" | cut -d':' -f2-)
      ;;
    "Format:")  # seen in /usr/share/doc/libpcsclite1/copyright
      result=$(grep '^License:' "$copyrightfile" | cut -d':' -f2-)
      ;;
    *)
      echo "WARNING: Unknown format of $copyrightfile: $format" >&2
      exit 1  # unknown format
  esac

  if [ -n "$result" ]; then
    echo "$result" | sed -r -e 's/ and /\n/g' -e 's/^ +//' -e 's/ +$//' -e 's/icence/icense/g' | sort -u
  fi

  exit 0
}

fuzzy_common_licenses () {
  set -e

  package="$1"
  copyrightfile=
  if [ -f "/usr/share/doc/$package/copyright" ]; then
    copyrightfile="/usr/share/doc/$package/copyright"
  elif [ -f "/usr/share/doc/${package%:*}/copyright" ]; then
    copyrightfile="/usr/share/doc/${package%:*}/copyright"
  else
    exit 0  # no copyright file found
  fi

  result=$(grep -oP '/usr/share/common-licenses/[0-9A-Za-z_.+-]+[0-9A-Za-z+]' "$copyrightfile" | cut -d/ -f5- | sort -u)
  if [ -n "$result" ]; then
    echo "$result"
  fi

  exit 0
}

superfuzzy () {
  set -e

  package="$1"
  copyrightfile=
  if [ -f "/usr/share/doc/$package/copyright" ]; then
    copyrightfile="/usr/share/doc/$package/copyright"
  elif [ -f "/usr/share/doc/${package%:*}/copyright" ]; then
    copyrightfile="/usr/share/doc/${package%:*}/copyright"
  else
    exit 0  # no copyright file found
  fi

  result=$(grep -Ewoi \
      -e '(4-?clause )?"?BSD"? licen[sc]es?' \
      -e '(Boost Software|mozilla (public)?|MIT) Licen[sc]es?' \
      -e '(CCPL|BSD|L?GPL)-[0-9a-z.+-]+( Licenses?)?' \
      -e 'Creative Commons( Licenses?)?' \
      -e 'Public Domain( Licenses?)?' \
      -e 'GNU General Public( License)?' \
      -e 'Info-ZIP License' \
      "$copyrightfile" | sed -r -e 's/[Ll]icence/License/g' | sort -u)
  if [ -n "$result" ]; then
    echo "$result"
  fi

  exit 0
}

free_licence () {
  set -e

  package="$1"
  copyrightfile=
  if [ -f "/usr/share/doc/$package/copyright" ]; then
    copyrightfile="/usr/share/doc/$package/copyright"
  elif [ -f "/usr/share/doc/${package%:*}/copyright" ]; then
    copyrightfile="/usr/share/doc/${package%:*}/copyright"
  else
    exit 0  # no copyright file found
  fi

  result=$(grep -e '^License:' -e '^Licence:' "$copyrightfile" | cut -d':' -f2-)
  if [ -n "$result" ]; then
    echo "$result" | sed -r -e 's/ and /\n/g' -e 's/^[[:space:]]*//' -e 's/[[:space:]]*$//' -e 's/icence/icense/g' | sort -u
  fi

  exit 0
}

dump_unreadable () {
  set -e

  package="$1"
  copyrightfile=
  if [ -f "/usr/share/doc/$package/copyright" ]; then
    copyrightfile="/usr/share/doc/$package/copyright"
  elif [ -f "/usr/share/doc/${package%:*}/copyright" ]; then
    copyrightfile="/usr/share/doc/${package%:*}/copyright"
  else
    exit 0  # no copyright file found
  fi

  #echo "$copyrightfile" >&2
  #head -n1 "$copyrightfile" >&2

  exit 0
}

methods=(machine_readable fuzzy_common_licenses superfuzzy free_licence dump_unreadable)

COLUMNS=2000 dpkg -l | grep '^.[iufhwt]' | while read pState package pVer pArch pDesc; do
  license=
  for method in $methods; do
    license=$("$method" "$package")
    [ $? -eq 0 ] || exit 1
    [ -n "$license" ] || continue
    # remove line breaks and spaces
    license=$(echo "$license" | tr '\n' ' ' | sed -r -e 's/ +/ /g' -e 's/^ +//' -e 's/ +$//')
    [ -z "$license" ] || break
  done
  [ -n "$license" ] || license='unknown'

  printf "$package $license\n"
done
`
