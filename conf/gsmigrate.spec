%define debug_package %{nil}
%define pkg_name      gsmigrate
%define utility_name  gsmigrate

Summary: Database migrations utility
Name: gsmigrate
Version: %{_app_version_number}
Release: %{_app_version_build}
License: Proprietary
Group: Productivity/WEB DESK/Utilities

URL: https://github.com/webnice/migrate
Source0: gsmigrate

PreReq:         %fillup_prereq /bin/mkdir /bin/cp
PreReq:         permissions
Requires:       bash
Requires(pre):  /usr/bin/chown /usr/bin/chmod
#Requires(post):
#Requires(preun):
#Requires(postun):

Provides: %{name} = %{version}-%{release}

BuildRoot: %{_tmppath}/%{name}-%{version}-%{release}-root
ExclusiveArch: x86_64
Conflicts:     wdstreamer


%description
Утилита применения миграций базы данных
Поддерживаемые базы данных:
- mysql
- postgres
- cockroach
- sqlite3
- redshift
- clickhouse
- tidb


%prep


%build


%install
rm -rf %{buildroot}
## bin
install -p -D -m 0755 %{SOURCE0} %{buildroot}%{_sbindir}/%{utility_name}


%pre
exit 0


%preun


%post
# System permissions
%set_permissions %{_sbindir}/%{utility_name}


%postun
exit 0


%clean
rm -rf %{buildroot}


%files
%defattr(-,root,root,-)
%attr(755,root,root) %{_sbindir}/%{utility_name}
%{_sbindir}/%{utility_name}


%changelog
* Mon Apr 2 2018 Alexander V. Tsarev <a.tsarev@webdesk.ru> version: %{_app_version_number} build: %{_app_version_build}
- Initial RPM (go binary) release.
