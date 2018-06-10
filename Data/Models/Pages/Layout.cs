using Dwn.Data.Models.Authoring;
using System;
using System.Collections.Generic;
using System.Text;

namespace Dwn.Data.Models.Pages
{
    public class Layout
    {
        public Guid Id { get; set; }
        public bool Vertical { get; set; }
        public List<AuthoredContent> AuthoredContents { get; set; }
    }
}
