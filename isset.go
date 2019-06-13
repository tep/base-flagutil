// Copyright © 2017 Tim Peoples <coders@toolman.org>
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package flagutil

/*
import "github.com/spf13/pflag"

func ValueIsSet(name string) (string, bool, error) {
	if err := MergeAndParse(); err != nil {
		return "", false, err
	}

	var isSet bool
	fs := pflag.Lookup(name)
	if fs == nil {
		return "", false, nil
	}

	value := fs.DefValue

	pflag.Visit(func(pf *pflag.Flag) {
		if pf.Name == name {
			isSet = true
			value = pf.Value.String()
		}
	})

	return value, isSet, nil
}
*/
