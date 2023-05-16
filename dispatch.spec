Name:           dispatch
Version:        0.5
Release:        1%{?dist}
Summary:        Provides an easy-to-use command to dispatch tasks described in a YAML file
License:        MIT

Source0:        %{name}-%{version}.tar.gz

BuildRequires:  golang
BuildRequires:  systemd-rpm-macros

Provides:       %{name} = %{version}

%description
Provides an easy-to-use command to dispatch tasks described in a YAML file

%global debug_package %{nil}

%prep
%autosetup

%build
go build -v -o %{name}

%install
mkdir -p %{buildroot}/usr/bin/
install -m 755 %{name} %{buildroot}/usr/bin/%{name}

%check

%post

%preun

%files
/usr/bin/%{name}

%changelog
