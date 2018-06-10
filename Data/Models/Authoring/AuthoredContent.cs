using Dwn.Data.Models.Identity;
using System;
using System.Collections.Generic;
using System.Text;

namespace Dwn.Data.Models.Authoring
{
    public class AuthoredContent
    {
        public Guid Id { get; set; }

        public string Text { get; set; }
        public ApplicationUser Author { get; set; }
        public DateTime Created { get; set; }
    }
}
