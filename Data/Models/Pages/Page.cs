using System;
using System.Collections.Generic;
using System.Text;

namespace Dwn.Data.Models.Pages
{
    public class Page
    {
        public Guid Id { get; set; }
        public string Slug { get; set; }
        public Layout Layout { get; set; }
    }
}
