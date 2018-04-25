using System;
using System.Collections.Generic;
using System.Text;

namespace Dwn.Cli
{
    public static class Helpers
    {
        public static T GetElementOrDefault<T>(this T[] array, int index)
        {
            if (index < array.Length)
            {
                return array[index];
            }
            return default(T);
        }
    }
}
