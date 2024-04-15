> [!IMPORTANT]
> *Deprecation notice*
>
> Recent versions of the Icinga Powershell Framework and Icinga2 mitigate the performance impact
> of checks. This helper is therefor not needed anymore and there will be no further development.


# icinga-powershell-connector

This is an **experimental** mitigation helper, trying to help with issues of high CPU usage with the [Icinga Powershell Framework] on Windows.

In this specific case, you already want to have installed the Icinga Powershell Service, with the
[REST API](https://icinga.com/docs/icinga-for-windows/latest/restapi/doc/02-Installation/) and the
[API checks](https://icinga.com/docs/icinga-for-windows/latest/apichecks/doc/01-Introduction/) enabled.
So checks can be run via API, and no Powershell has to be executed for the checks.

The edge case, the connector want to fix, is to remove the requirement to starting a Powershell for retrieving data
from the REST API. With the normal CheckCommands, Icinga 2 will start the Powershell and withing functions of the 

**Note:** Please be aware, that Icinga is trying to fix this issue soon, and it usually should only cause problems on
smaller systems or VMs, with some Anti Virus software auditing the startup of Powershell.

[Icinga Powershell Framework]: https://icinga.com/docs/icinga-for-windows/latest

**Further reading:**

* https://github.com/Icinga/icinga-powershell-framework/issues/131
* https://github.com/Icinga/icinga2/issues/8082

Internal notes: ref/NC/730122

## Requirements

* Icinga 2 as agent on Windows
* Powershell Framework >= 1.5 with restapi and apichecks enabled
* The apichecks REST API service listening on `https://localhost:5668`
* Director or Icinga config files CheckCommands for the Powershell Framework ready to be changed

## Installation

To install the connector, download it from the [releases] section on GitHub and copy the file onto the Windows systems.

We would recommend: `C:\Program Files\Icinga2\sbin\powershell-connector.exe`

After you made sure to install the connector onto all your Windows systems, you can change the `PowerShell Base` command
in Icinga 2, so it will use the connector.

The new command will be the location you installed the connector at:

* `C:\Program Files\Icinga2\sbin\powershell-connector.exe`
* `PluginDir + "\powershell-connector.exe"` (for plain config)
* `powershell-connector.exe` or `PluginDir + powershell-connector.exe` (for Director)

[releases]: https://github.com/NETWAYS/icinga-powershell-connector/releases

## Example

Usually for the Powershell checks are executed from Icinga like this:

```
'C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe' `
    '-C' "try { Use-Icinga -Minimal; } catch { Write-Output 'some error message'; exit 3; }; Exit-IcingaExecutePlugin -Command 'Invoke-IcingaCheckUsedPartitionSpace' " `
    '-Warning' '80' '-Critical' '95' '-Include' '@()' '-Exclude' '@()' '-Verbosity' '2'
```

The idea with the connector is, to just replace the exe, and keep the rest identical. URL and certificates are handled
automatically.

```
'C:\Program Files\Icinga2\sbin\powershell-connector.exe' `
    '-C' "try { Use-Icinga -Minimal; } catch { Write-Output 'some error message'; exit 3; }; Exit-IcingaExecutePlugin -Command 'Invoke-IcingaCheckUsedPartitionSpace' " `
    '-Warning' '80' '-Critical' '95' '-Include' '@()' '-Exclude' '@()' '-Verbosity' '2'
```

## License

Copyright (C) 2021 [NETWAYS GmbH](mailto:info@netways.de)

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.
If not, see [gnu.org/licenses](https://www.gnu.org/licenses/).
