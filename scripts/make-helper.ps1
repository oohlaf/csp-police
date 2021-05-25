<#
.Synopsis
   Helper functions that can be called from a Makefile on Windows.
.DESCRIPTION
   Helper methods for:

   help           Show usage information on all targets in a Makefile.
   rm             Remove files.
   rmdir          Remove files or directories
.EXAMPLE
   PowerShell -ExecutionPolicy ByPass -File ./scripts/make-helper.ps1 rm file1 file2
.EXAMPLE
   PowerShell -ExecutionPolicy ByPass -File ./scripts/make-helper.ps1 help $(MAKEFILE_LIST)
.INPUTS
   Name of the function (help, rm).
   Optional arguments to that function.
.NOTES
  Version:        1.0
  Author:         Olaf Conradi <olaf@conradi.org>
#>
[CmdletBinding()]
Param(
   [Parameter(Mandatory=$true, Position=0)]
   [ValidateNotNull()]
   [ValidateSet("help", "rm", "rmdir")]
   $Func,

   [Parameter(Mandatory=$false, Position=1, ValueFromRemainingArguments=$true)]
   [string[]]
   $Arguments
)

<#
.Synopsis
   Show overview of target descriptions within a Makefile.
.DESCRIPTION
   A helper to show usage information.
   It will extract targets with descriptions starting with `##- `.
   Lines starting with # are ignored.
.EXAMPLE
   Invoke-Help Makefile
.INPUTS
   The contents of a Makefile.
.OUTPUTS
   The usage overview of all targets and descriptions.
#>
function Invoke-Help {
   Param(
      [Parameter(Mandatory=$true)]
      [string[]]
      $Makefile
   )
   Write-Host "Usage: make [TARGET]"
   Write-Host ""
   Write-Host "Targets:"
   Get-Content $Makefile | Select-String -Pattern "^([^#].*):.*##-\s(.*)$" | ForEach-Object {
      "  {1,-28}{2}" -f $_.Matches.Groups
   } | Write-Host
}

<#
.Synopsis
   Remove files
.DESCRIPTION
   A helper to call Remove-Item from a Makefile.
   
   Within a Makefile all paths follow Unix directory separators which is not understood by
   the delete command within cmd.exe. Powershell does support it.
.EXAMPLE
   Invoke-Remove -Path folder/file.txt folder/file2.txt
.INPUTS
   Path is an array of files to remove.
#>
function Invoke-Remove {
   Param(
      [Parameter(Mandatory=$true)]
      [string[]]
      $Path
   )
   Remove-Item -ErrorAction SilentlyContinue -Force -Path $Path
}

<#
.Synopsis
   Remove directories
.DESCRIPTION
   A helper to call Remove-Item with option recursive from a Makefile.
   
   Within a Makefile all paths follow Unix directory separators which is not understood by
   the delete command within cmd.exe. Powershell does support it.
.EXAMPLE
   Invoke-RemoveDir -Path folder/file.txt folder/file2.txt
.INPUTS
   Path is an array of files or directories to remove.
#>
function Invoke-RemoveDir {
   Param(
      [Parameter(Mandatory=$true)]
      [string[]]
      $Path
   )
   Remove-Item -ErrorAction SilentlyContinue -Force -Recurse -Path $Path
}

switch ($Func) {
   'help' {
      Invoke-Help -Makefile $Arguments
   }
   'rm' {
      Invoke-Remove -Path $Arguments
   }
   'rmdir' {
      Invoke-RemoveDir -Path $Arguments
   }
   default {
      Write-Error "Unknown helper function"
   }
}
