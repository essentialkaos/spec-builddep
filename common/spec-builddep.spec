################################################################################

%define debug_package  %{nil}

################################################################################

Summary:        Utility for installing dependencies for building an RPM package
Name:           spec-builddep
Version:        1.1.1
Release:        0%{?dist}
Group:          Applications/System
License:        Apache License, Version 2.0
URL:            https://kaos.sh/spec-builddep

Source0:        https://source.kaos.st/%{name}/%{name}-%{version}.tar.bz2

BuildRoot:      %{_tmppath}/%{name}-%{version}-%{release}-root-%(%{__id_u} -n)

BuildRequires:  golang >= 1.23

Requires:       rpm rpm-build

Provides:       %{name} = %{version}-%{release}

################################################################################

%description
spec-builddep is a simple utility for installing dependencies for building
an RPM package (yum-builddep drop-in replacement).

################################################################################

%prep

%setup -q
if [[ ! -d "%{name}/vendor" ]] ; then
  echo -e "----\nThis package requires vendored dependencies\n----"
  exit 1
elif [[ -f "%{name}/%{name}" ]] ; then
  echo -e "----\nSources must not contain precompiled binaries\n----"
  exit 1
fi

%build
pushd %{name}
  %{make_build} all
  cp LICENSE ..
popd

%install
rm -rf %{buildroot}

install -dm 755 %{buildroot}%{_bindir}
install -dm 755 %{buildroot}%{_mandir}/man1

install -pm 755 %{name}/%{name} %{buildroot}%{_bindir}/

./%{name}/%{name} --generate-man > %{buildroot}%{_mandir}/man1/%{name}.1

%clean
rm -rf %{buildroot}

%post
if [[ -d %{_sysconfdir}/bash_completion.d ]] ; then
  %{name} --completion=bash 1> %{_sysconfdir}/bash_completion.d/%{name} 2>/dev/null
fi

if [[ -d %{_datarootdir}/fish/vendor_completions.d ]] ; then
  %{name} --completion=fish 1> %{_datarootdir}/fish/vendor_completions.d/%{name}.fish 2>/dev/null
fi

if [[ -d %{_datadir}/zsh/site-functions ]] ; then
  %{name} --completion=zsh 1> %{_datadir}/zsh/site-functions/_%{name} 2>/dev/null
fi

%postun
if [[ $1 == 0 ]] ; then
  if [[ -f %{_sysconfdir}/bash_completion.d/%{name} ]] ; then
    rm -f %{_sysconfdir}/bash_completion.d/%{name} &>/dev/null || :
  fi

  if [[ -f %{_datarootdir}/fish/vendor_completions.d/%{name}.fish ]] ; then
    rm -f %{_datarootdir}/fish/vendor_completions.d/%{name}.fish &>/dev/null || :
  fi

  if [[ -f %{_datadir}/zsh/site-functions/_%{name} ]] ; then
    rm -f %{_datadir}/zsh/site-functions/_%{name} &>/dev/null || :
  fi
fi

################################################################################

%files
%defattr(-,root,root,-)
%doc LICENSE
%{_mandir}/man1/%{name}.1.*
%{_bindir}/%{name}

################################################################################

%changelog
* Thu May 15 2025 Anton Novojilov <andy@essentialkaos.com> - 1.1.1-0
- Code refactoring
- Dependencies update

* Thu Apr 17 2025 Anton Novojilov <andy@essentialkaos.com> - 1.1.0-0
- Improved quiet mode
- Updated usage info
- Fix error output in quiet mode
- Dependencies update

* Tue Sep 10 2024 Anton Novojilov <andy@essentialkaos.com> - 1.0.3-0
- Fixed bug with handling defined macro

* Fri Sep 06 2024 Anton Novojilov <andy@essentialkaos.com> - 1.0.2-0
- Added support of conditions with typos ('=>', '=<')
- Package ek updated to v13
- Code refactoring
- Dependencies update

* Mon Jun 24 2024 Anton Novojilov <andy@essentialkaos.com> - 1.0.1-0
- Code refactoring
- Dependencies update

* Thu Mar 28 2024 Anton Novojilov <andy@essentialkaos.com> - 1.0.0-0
- Improved support information gathering
- Code refactoring
- Dependencies update

* Mon Oct 09 2023 Anton Novojilov <andy@essentialkaos.com> - 0.1.2-0
- Better check for installed packages

* Fri Oct 06 2023 Anton Novojilov <andy@essentialkaos.com> - 0.1.1-0
- Improve check for installed versions on EL 7

* Sun Oct 01 2023 Anton Novojilov <andy@essentialkaos.com> - 0.1.0-0
- Resolve real packages names before check

* Mon Sep 25 2023 Anton Novojilov <andy@essentialkaos.com> - 0.0.3-0
- Added dependencies deduplication
- Fixed check for rpm-build
- Added exit code of yum/dnf to error message about failed install

* Mon Sep 25 2023 Anton Novojilov <andy@essentialkaos.com> - 0.0.2-0
- Add '--exclude' option support

* Sun Sep 24 2023 Anton Novojilov <andy@essentialkaos.com> - 0.0.1-0
- Initial release
