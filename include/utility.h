/* 
    netcalc - a simple CLI subnet calculator written in ANSI C
    utility.h - misc function definitions 

    Copyright (C) 2009  Matthew Reath 

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.
    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

#ifndef __UTILITY_H__
#define __UTILITY_H__

/* dotted decimal string to integer */
unsigned int ddtoint(char* dd);

/* integer to dotted decimal */
int inttodd(char* dd, unsigned int ip);

#endif
