#! /usr/bin/python3
# pylint: disable=too-few-public-methods
# pylint: disable=no-self-use
'''
Checks versioning command
'''

import subprocess
import sys
from modules.Messages import Message


class VersionTest():
    ''' Nsulate Version Test '''

    def test(self, git_arg, version_arg):
        ''' Runs the test '''

        print(git_arg)
        print(version_arg)

        version = subprocess.getoutput("nsulate version")
        subprocess.getoutput("nsulated -v 2> tmp.txt")
        daemon_version = subprocess.getoutput(
            "sed -n 1p tmp.txt | awk '{print $3}'")

        if version == daemon_version and version != 0:
            Message.status("Nsulate versions match.")
        else:
            Message.error(
                "Versions do not match. Nsulate: " +
                version +
                " Nsulated: " +
                daemon_version)

        # Will not work until code has been written to give script the correct
        # values

        # subprocess.Popen("nsulate version -a >> tmp2.txt", shell=True)
        # nsulate_git = subprocess.getoutput("sed -n 2p tmp2.txt | awk '{print $4}'")
        #
        # if (git_arg == nsulate_git):
        #     Message.status("Nsulate is using correct git commit.")
        # else:
        #     Message.error("Commits do not match. Nsulate: " + git + " Local: " + nsulate_git)

        #
        # if( version_arg == version ):
        #   Message.status("Nsulate versions match.")
        # else:
        #   Message.error("Versions do not match. Nsulate: "
        # + version + " Nsulated: " + daemon_version)


def main():
    ''' The main method '''

    # python VersionTest 17.2.4442 11111

    git_arg = sys.argv[0]
    version_arg = sys.argv[1]

    # Create and run test
    versiontest = VersionTest()
    versiontest.test(git_arg, version_arg)


if __name__ == "__main__":
    main()
