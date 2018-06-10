using System;
using System.Collections.Generic;
using System.Text;
using Dwn.Data.Models.Authoring;

namespace Dwn.Data.Models.Blogging
{
    public class Post : AuthoredContent
    {
        public string Title { get; set; }
        public string Slug { get; set; }
        public BlogCategory BlogCategory { get; set; }

    }
}
