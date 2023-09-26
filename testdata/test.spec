################################################################################

%{!?_without_check: %define _with_check 1}

%define pg_ver  10
%define pg_fullver  %{pg_ver}.23

################################################################################

Summary:            Test spec
Name:               myapp
Version:            1.0.0
Release:            0%{?dist}
Group:              System Environment/Base
License:            MIT
URL:                https://domain.com

BuildRequires:      postgresql%{pg_ver}-devel = %{pg_fullver}
BuildRequires:      make gcc zlib-devel >= 1.2.11 readline-devel <= 7
BuildRequires:      python < 3 perl > 5 bash = 4 gcc make
BuildRequires:      perl(ExtUtils::Embed) python3dist(setuptools)
# BuildRequires:    gdal-devel gcc-c++ >= 8

BuildRoot:          %{_tmppath}/%{name}-%{version}-%{release}-root-%(%{__id_u} -n)

Source0:            https://source.kaos.st/perfecto/%{name}-%{version}.tar.gz

################################################################################

%description
Test spec.

################################################################################

%package magic

Summary:            Test subpackage for perfecto
Group:              System Environment/Base

%description magic
Test subpackage.

################################################################################

%prep
%setup -qn %{name}-%{version}

%build
%{__make} %{?_smp_mflags}

%install
rm -rf %{buildroot}

%{make_install} PREFIX=%{buildroot}%{_prefix}

%clean
rm -rf %{buildroot}

%check
%if %{?_with_check:1}%{?_without_check:0}
%{make} check
%endif

%post
%{__chkconfig} --add %{name} &>/dev/null || :

%preun
%{__chkconfig} --del %{name} &> /dev/null || :

%postun
%{__chkconfig} --del %{name} &> /dev/null || :

################################################################################

%files
%defattr(-,root,root,-)
%{_bindir}/%{name}

%files magic
%defattr(-,root,root,-)
%{_bindir}/%{name}-magic

################################################################################

%changelog
* Wed Jan 24 2018 Anton Novojilov <andy@essentialkaos.com> - 1.0.0-0
- Test changelog record
