Name:           dispatch
Version:        0.3
Release:        1%{?dist}
Summary:        Provides an easy-to-use command to dispatch tasks described in a YAML file

License:        GPLv3
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

%check

%post

%preun

%files

%changelog
