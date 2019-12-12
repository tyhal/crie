#include "example.h"
#include <algorithm>
#include <array>
#include <functional>
#include <iterator>

int test::myClass::afunction( int a )
{
    //	int unused = 0;
    return a;
}

int main()
{
    std::array< int, 5 > arr = { 3, 4, 1, 5, 2 };
    std::sort( std::begin( arr ), std::end( arr ) );
    std::sort( std::begin( arr ), std::end( arr ), std::greater< int >{} );
    test::myClass myclass;
    myclass.afunction( 0 );
    return 0;
}
